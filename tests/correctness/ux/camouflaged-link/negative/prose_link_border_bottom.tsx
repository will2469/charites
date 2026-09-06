export function ParagraphLinkBorderBottom() {
  return (
    <article className="prose">
      <p>
        Untuk informasi lebih lanjut, silakan hubungi{" "}
        <a href="mailto:info@example.com" className="text-primary font-medium border-b border-primary">
          Layanan Pelanggan
        </a>{" "}
        kami.
      </p>
    </article>
  );
}
