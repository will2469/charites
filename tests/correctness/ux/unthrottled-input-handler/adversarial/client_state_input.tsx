// Adversarial fixture: Input updating local component state without network calls
export function ControlledInput() {
  const [val, setVal] = useState("");

  return (
    <input
      type="text"
      value={val}
      onChange={(e) => setVal(e.target.value)}
      className="border rounded px-3 py-2"
    />
  );
}
