export function SvgButton() {
  return (
    <button className="flex items-center gap-2">
      <svg className="size-6"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z" /></svg>
      <span>Beranda</span>
    </button>
  );
}

export function SpreadIcon(props: any) {
  return (
    <button>
      <span className="material-icons" {...props}>favorite</span>
    </button>
  );
}
