export function CompliantFormControls() {
  return (
    <div>
      <FormItem>
        <FormLabel>Nama Depan</FormLabel>
        <FormControl>
          <input className="h-11 px-3.5 py-2.5 border" placeholder="John" />
        </FormControl>
      </FormItem>
      <FormItem>
        <FormLabel>Periode</FormLabel>
        <FormControl>
          <fieldset>
            <legend className="sr-only">Periode Libur</legend>
            <DateRangePicker />
          </fieldset>
        </FormControl>
      </FormItem>
    </div>
  );
}
