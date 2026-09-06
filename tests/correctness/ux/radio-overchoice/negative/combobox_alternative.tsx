import React from 'react';

const cityOptions = [
  { value: 'jkt', label: 'Jakarta' },
  { value: 'bdg', label: 'Bandung' },
  { value: 'sby', label: 'Surabaya' },
  { value: 'smg', label: 'Semarang' },
  { value: 'mdn', label: 'Medan' },
  { value: 'mks', label: 'Makassar' },
  { value: 'dps', label: 'Denpasar' },
  { value: 'plm', label: 'Palembang' },
];

export function SearchableCitySelect() {
  return (
    <div className="p-4 space-y-2">
      <label className="text-sm font-semibold">Pilih Kota Tujuan</label>
      <Combobox
        options={cityOptions}
        placeholder="Cari kota..."
        searchable
      />
    </div>
  );
}
