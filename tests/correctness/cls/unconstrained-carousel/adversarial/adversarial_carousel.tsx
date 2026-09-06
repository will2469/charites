export function DynamicCarousel(props: any) {
  return <Carousel {...props} />;
}

export function ScrollableParagraph() {
  return (
    <div className="overflow-x-auto p-4">
      <p>Teks deskripsi biasa tanpa fitur carousel snapping.</p>
    </div>
  );
}
