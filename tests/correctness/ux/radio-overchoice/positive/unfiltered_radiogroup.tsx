import React from 'react';

export function UnfilteredCityPicker() {
  return (
    <div className="p-4">
      <RadioGroup name="city" className="space-y-2">
        <RadioGroupItem value="jkt" label="Jakarta" />
        <RadioGroupItem value="bdg" label="Bandung" />
        <RadioGroupItem value="sby" label="Surabaya" />
        <RadioGroupItem value="smg" label="Semarang" />
        <RadioGroupItem value="mdn" label="Medan" />
        <RadioGroupItem value="mks" label="Makassar" />
        <RadioGroupItem value="dps" label="Denpasar" />
        <RadioGroupItem value="plm" label="Palembang" />
      </RadioGroup>
    </div>
  );
}
