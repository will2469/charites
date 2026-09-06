export function NestedScrollContainer() {
  return (
    <div className="overflow-auto max-w-full">
      <div className="p-2 border rounded">
        <table className="w-full">
          <tbody>
            <tr>
              <td>Nested table inside ancestor overflow-auto</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
}
