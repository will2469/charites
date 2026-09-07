package ux_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/ux"
)

func TestDestructiveActionUnconfirmedRule(t *testing.T) {
	rule := ux.NewDestructiveActionUnconfirmedRule()

	cases := []struct {
		name      string
		code      string
		wantDiags int
	}{
		{
			name: "Case 1: Direct deleteUser call without confirmation",
			code: `
export function DeleteDirect() {
  return (
    <button
      onClick={() => deleteUser(user.id)}
      className="bg-destructive text-destructive-foreground"
    >
      Hapus Pengguna
    </button>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 2: Wrapped in AlertDialogTrigger",
			code: `
export function DeleteWithTrigger() {
  return (
    <AlertDialogTrigger asChild>
      <button className="bg-destructive text-destructive-foreground">
        Hapus Pengguna
      </button>
    </AlertDialogTrigger>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 3: Inline window.confirm check",
			code: `
export function DeleteWithWindowConfirm() {
  return (
    <button
      onClick={() => {
        if (window.confirm("Apakah Anda yakin ingin menghapus?")) {
          deleteUser(id);
        }
      }}
      className="bg-destructive"
    >
      Hapus Pengguna
    </button>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 4: Staged setConfirmIndex -> open={confirmIndex !== null} -> onConfirm (Issue #11 specimen)",
			code: `
export function FollowerList({ users, onHapus }: Props) {
  const [confirmIndex, setConfirmIndex] = useState<number | null>(null);

  return (
    <div>
      {users.map((user, index) => (
        <Button
          key={user.id}
          type="button"
          variant="destructive"
          onClick={() => setConfirmIndex(index)}
        >
          Hapus
        </Button>
      ))}

      <ActionApprovalDialog
        open={confirmIndex !== null}
        title="Hapus Pengikut"
        onConfirm={() => {
          if (confirmIndex !== null) onHapus(confirmIndex);
          setConfirmIndex(null);
        }}
        onCancel={() => setConfirmIndex(null)}
      />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 5: setIsOpen(true) with informational dialog (no confirmation action)",
			code: `
export function InfoModalOnly() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div>
      <Button variant="destructive" onClick={() => setIsOpen(true)}>
        Hapus
      </Button>

      <Dialog open={isOpen}>
        <div>Informasi penghapusan pengguna...</div>
      </Dialog>
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 6: Unrelated setSelectedTab() with direct deleteUser in same handler",
			code: `
export function TabAndDirectDelete() {
  const [selectedTab, setSelectedTab] = useState("all");

  return (
    <Button
      variant="destructive"
      onClick={() => {
        setSelectedTab("all");
        deleteUser(id);
      }}
    >
      Hapus
    </Button>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 7: Setter and dialog unrelated to state",
			code: `
export function UnrelatedStateDialog() {
  const [selectedTab, setSelectedTab] = useState("all");
  const [otherState, setOtherState] = useState(false);

  return (
    <div>
      <Button variant="destructive" onClick={() => setSelectedTab("archived")}>
        Hapus
      </Button>
      <ConfirmDialog open={otherState} onConfirm={() => doSomething()} />
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 8: State declared in another component (out of component scope)",
			code: `
function Parent() {
  const [confirmOpen, setConfirmOpen] = useState(false);
  return (
    <>
      <Child />
      <ConfirmDialog open={confirmOpen} onConfirm={() => doDelete()} />
    </>
  );
}

function Child() {
  return (
    <Button variant="destructive" onClick={() => setConfirmOpen(true)}>
      Hapus
    </Button>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 9: Setter passed as prop",
			code: `
export function ChildComponent({ setConfirmOpen }: Props) {
  return (
    <Button variant="destructive" onClick={() => setConfirmOpen(true)}>
      Hapus
    </Button>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 10: setDeleteDialogOpen(true) with valid downstream confirm action",
			code: `
export function DeleteModalSafe() {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  return (
    <div>
      <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
        Hapus
      </Button>

      <AlertDialog open={deleteDialogOpen} onConfirm={() => executePurge()}>
        <AlertDialogContent>
          <AlertDialogTitle>Konfirmasi</AlertDialogTitle>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 11: setPendingDeleteId(id) -> open={Boolean(pendingDeleteId)} -> onConfirm",
			code: `
export function BooleanPredicateSafe() {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);

  return (
    <div>
      <Button variant="destructive" onClick={() => setPendingDeleteId(user.id)}>
        Hapus
      </Button>

      <ConfirmDialog
        open={Boolean(pendingDeleteId)}
        onConfirm={() => removeUser(pendingDeleteId)}
      />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 12: pendingDeleteId && <ConfirmDialog onConfirm={...}> (conditional render guard)",
			code: `
export function ConditionalGuardSafe() {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);

  return (
    <div>
      <Button variant="destructive" onClick={() => setPendingDeleteId(user.id)}>
        Hapus
      </Button>

      {pendingDeleteId && (
        <ConfirmDialog onConfirm={() => removeUser(pendingDeleteId)} />
      )}
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 13: pendingDeleteId !== null but informational Dialog",
			code: `
export function InfoDialogWithId() {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);

  return (
    <div>
      <Button variant="destructive" onClick={() => setPendingDeleteId(user.id)}>
        Hapus
      </Button>

      <Dialog open={pendingDeleteId !== null}>
        <p>Pending operation...</p>
      </Dialog>
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 14: Direct delete + setConfirmOpen(true) in same handler",
			code: `
export function MixedDirectDeleteAndSetter() {
  const [confirmOpen, setConfirmOpen] = useState(false);

  return (
    <div>
      <Button
        variant="destructive"
        onClick={() => {
          setConfirmOpen(true);
          deleteUser(id);
        }}
      >
        Hapus
      </Button>

      <ConfirmDialog open={confirmOpen} onConfirm={() => deleteUser(id)} />
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 15: Computed/dynamic state expression not in supported predicates",
			code: `
export function ComputedStateExpression() {
  const [pendingCount, setPendingCount] = useState(0);

  return (
    <div>
      <Button variant="destructive" onClick={() => setPendingCount(1)}>
        Hapus
      </Button>

      <ConfirmDialog
        open={pendingCount > 0 && isReady}
        onConfirm={() => purgeAll()}
      />
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 16: Unknown setter naming resolved via useState binding",
			code: `
export function CustomSetterBindingSafe() {
  const [pendingAction, updatePendingAction] = useState<string | null>(null);

  return (
    <div>
      <Button
        variant="destructive"
        onClick={() => updatePendingAction(item.id)}
      >
        Hapus
      </Button>

      <ActionApprovalDialog
        open={pendingAction !== null}
        onConfirm={() => executeAction(pendingAction)}
      />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 17: Multiple destructive buttons mapping to same confirmation state (list/table)",
			code: `
export function TableMultipleButtons() {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);

  return (
    <div>
      <table>
        <tbody>
          <tr>
            <td>
              <Button
                variant="destructive"
                onClick={() => setPendingDeleteId("1")}
              >
                Hapus
              </Button>
            </td>
          </tr>
          <tr>
            <td>
              <Button
                variant="destructive"
                onClick={() => setPendingDeleteId("2")}
              >
                Hapus
              </Button>
            </td>
          </tr>
        </tbody>
      </table>

      <ConfirmDialog
        open={pendingDeleteId !== null}
        onConfirm={() => deleteItem(pendingDeleteId)}
      />
    </div>
  );
}`,
			wantDiags: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := tsx.Extract([]byte(tc.code))
			if err != nil {
				t.Fatalf("failed to parse TSX: %v", err)
			}

			var diagsCount int
			for node := range root.Walk() {
				findings := rule.Evaluate(node)
				diagsCount += len(findings)
			}

			if diagsCount != tc.wantDiags {
				t.Errorf("[%s] got %d diagnostics, want %d", tc.name, diagsCount, tc.wantDiags)
			}
		})
	}
}
