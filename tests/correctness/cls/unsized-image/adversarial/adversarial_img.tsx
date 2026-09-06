export function DynamicSpreadImage(props: any) {
  return (
    <div className="container">
      <img src="/photo.jpg" alt="Dynamic Photo" {...props} />
    </div>
  );
}

export function ParentAspectWrapper() {
  return (
    <div className="w-full aspect-video">
      <img src="/banner.jpg" alt="Banner" className="w-full h-full object-cover" />
    </div>
  );
}
