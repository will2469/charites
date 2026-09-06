export function DynamicSpreadIframe(props: any) {
  return (
    <div>
      <iframe src="/embed/chart" title="Chart" {...props} />
    </div>
  );
}

export function NestedGrandparentAspect() {
  return (
    <div className="aspect-video">
      <div className="relative w-full h-full">
        <iframe src="https://youtube.com/embed/123" title="Video" className="w-full h-full" />
      </div>
    </div>
  );
}
