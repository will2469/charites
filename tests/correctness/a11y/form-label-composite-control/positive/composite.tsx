export function CompositeFormControls() {
  return (
    <div>
      <FormItem>
        <FormLabel>Rentang Tanggal Libur</FormLabel>
        <FormControl>
          <DateRangePicker />
        </FormControl>
      </FormItem>
      <FormItem>
        <FormLabel>Alamat Lengkap</FormLabel>
        <FormControl>
          <AddressFields />
        </FormControl>
      </FormItem>
    </div>
  );
}
