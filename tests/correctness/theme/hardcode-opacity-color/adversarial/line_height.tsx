import React from "react";

export function ArticleTypography() {
  return (
    <article className="prose">
      <h1 className="text-3xl/9 font-extrabold text-foreground">Article Title</h1>
      <p className="text-sm/6 text-muted font-normal">
        Subtitle with fractional line-height modifier text-sm/6.
      </p>
      <p className="text-xs/relaxed text-muted">
        Another paragraph with text-xs/relaxed and text-base/7.
      </p>
    </article>
  );
}
