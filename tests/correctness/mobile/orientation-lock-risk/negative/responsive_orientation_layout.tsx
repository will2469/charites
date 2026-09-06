import React from "react";

export function ResponsiveOrientationLayout() {
  return (
    // Compliant: fluid responsive layout adapting naturally via landscape CSS variant
    <div className="flex flex-col landscape:flex-row gap-4 p-4">
      <aside className="w-full landscape:w-64">Navigasi</aside>
      <main className="flex-1">Konten Utama</main>
    </div>
  );
}
