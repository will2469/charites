// NEG-002: Other meta tags without theme-color
export function ValidMeta() {
  return (
    <head>
      <meta charSet="utf-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <meta name="description" content="Awesome application" />
    </head>
  );
}
