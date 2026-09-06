export function CompliantFormItem() {
  return (
    <div>
      <FormItem>
        <FormLabel>Nama Lengkap</FormLabel>
        <FormControl>
          <input className="h-11 px-3.5 py-2.5 border" placeholder="Nama Anda" />
        </FormControl>
      </FormItem>
      <FormItem>
        <FormLabel>Email</FormLabel>
        <input type="email" className="h-11 px-3.5 py-2.5 border" placeholder="Email" />
      </FormItem>
    </div>
  );
}
