export function SpreadAd(props: any) {
  return (
    <div>
      <div id="ad-dynamic" data-ad-slot="999" {...props} />
    </div>
  );
}

export function AddressBox() {
  return (
    <div id="address-box" className="p-4 bg-card rounded-lg">
      <p>Kantor Pelayanan Desa Sukamaju</p>
    </div>
  );
}
