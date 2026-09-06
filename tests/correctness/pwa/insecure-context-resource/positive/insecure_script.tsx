import React from "react";

export function InsecureScriptAndLink() {
  return (
    <div>
      {/* Insecure HTTP script violates Secure Contexts */}
      <script src="http://cdn.tracker.com/script.js"></script>
    </div>
  );
}
