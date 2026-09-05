// NEG-001: Using unified useTheme hook
import { useTheme } from "./theme-provider";

export function SafeToggle() {
  const { theme, setTheme } = useTheme();
  return (
    <button onClick={() => setTheme(theme === "dark" ? "light" : "dark")}>
      Toggle Theme
    </button>
  );
}
