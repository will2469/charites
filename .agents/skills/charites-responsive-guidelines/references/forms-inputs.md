# Mobile Forms & Inputs -- Deep Dive

Open this file when building any form with NIK, phone, email, date, address, file upload, camera capture, OTP, or multi-step flows.

## Table of Contents

1. [The Four-Attribute Combo](#1-the-four-attribute-combo)
2. [Indonesian-Specific Inputs (NIK, KK, NIP)](#2-indonesian-specific-inputs)
3. [Phone Numbers](#3-phone-numbers)
4. [Email](#4-email)
5. [Dates & Times](#5-dates--times)
6. [Address & Geographic](#6-address--geographic)
7. [File Upload & Camera](#7-file-upload--camera)
8. [OTP / Verification Codes](#8-otp-verification-codes)
9. [Multi-Step Forms](#9-multi-step-forms)
10. [Validation & Error Display](#10-validation--error-display)

---

## 1. The Four-Attribute Combo

Every mobile form input should set four attributes in concert:

| Attribute      | What it does                                                                                                |
| :------------- | :---------------------------------------------------------------------------------------------------------- |
| `type`         | Determines validation, native picker, and base keyboard                                                     |
| `inputMode`    | Hints the keyboard layout (numeric, tel, email, etc.) -- overrides `type` for keyboard only, not validation |
| `autoComplete` | Maps to the browser's autofill database; users get one-tap fill                                             |
| `enterKeyHint` | Labels the keyboard's submit key ("Next", "Go", "Done", etc.)                                               |

**Why all four**: `type="tel"` gives you tel validation but a telephone keypad. `type="text" inputMode="tel"` gives you the keypad without enforcing telephone-format validation. `autoComplete="tel"` lets the user autofill. `enterKeyHint="next"` moves to the next field. Together they create a native-feel flow.

---

## 2. Indonesian-Specific Inputs

### NIK (Nomor Induk Kependudukan) -- 16 digits

```tsx
<input
	type="text"
	inputMode="numeric"
	pattern="[0-9]*"
	maxLength={16}
	minLength={16}
	autoComplete="off" // NIK is sensitive -- never autofill
	enterKeyHint="next"
	placeholder="16 digit NIK"
	inputClassName="tracking-wider font-mono" // monospace for digit alignment
/>
```

**Notes**:

- `type="text"` not `type="number"` -- `number` drops leading zeros, which NIK can have.
- `pattern="[0-9]*"` triggers the numeric keyboard on iOS Safari (legacy).
- `autoComplete="off"` -- NIK is personally identifying; don't store in autofill.
- `font-mono` + `tracking-wider` makes 16 digits easier to read and compare against a physical KTP.
- Format display: some teams mask middle digits `3201**********1234` for privacy in lists, but the input should show full digits.

### NKK (Nomor Kartu Keluarga) -- 16 digits

Same as NIK. `autoComplete="off"`.

### NIP (Nomor Induk Pegawai) -- 18 digits

Same pattern, 18 digits.

### NPWP (Nomor Pokok Wajib Pajak) -- 15 digits, formatted `XX.XXX.XXX.X-XXX.XXX`

```tsx
const formatNpwp = (raw: string) => {
	const d = raw.replace(/\D/g, "").slice(0, 15);
	return d
		.replace(/^(\d{2})(\d)/, "$1.$2")
		.replace(/^(\d{2})\.(\d{3})(\d)/, "$1.$2.$3")
		.replace(/\.(\d{3})(\d)/, ".$1-$2")
		.replace(/-(\d{3})(\d)/, "-$1.$2");
};

<input
	type="text"
	inputMode="numeric"
	value={formatNpwp(value)}
	onChange={(e) => setValue(e.target.value.replace(/\D/g, ""))}
	autoComplete="off"
	placeholder="00.000.000.0-000.000"
/>;
```

### License plates (Plat Nomor)

`B 1234 ABC` format. Use `type="text" inputMode="text"` with `autoCapitalize="characters"`.

---

## 3. Phone Numbers

### Indonesian phone format

```tsx
<input
	type="tel"
	inputMode="tel"
	autoComplete="tel"
	autoCapitalize="none"
	pattern="(\+62|62|0)8[1-9][0-9]{6,11}"
	maxLength={15}
	enterKeyHint="next"
	placeholder="0812-3456-7890"
/>
```

**Notes**:

- `type="tel"` -- no strict validation, but gives tel keypad.
- `pattern` -- Indonesian mobile numbers start with `08`, `+628`, or `628`.
- Don't force `+62` -- Indonesian users overwhelmingly use the `08` format.
- Format display: `0812-3456-7890` is more readable than `08123456789`. Use a formatter in `onChange`.

### Why not `type="number"` for phone

`type="number"` is for quantities, not phone numbers. It strips leading `+`, drops leading zeros, shows spinner arrows, and doesn't accept formatting characters. Always `type="tel"`.

---

## 4. Email

```tsx
<input
	type="email"
	inputMode="email"
	autoComplete="email"
	autoCapitalize="none"
	autoCorrect="off"
	spellCheck={false}
	enterKeyHint="next"
/>
```

**`autoCapitalize="none"`** is critical -- without it, iOS auto-capitalizes the first letter of every email, which is wrong.

**`autoCorrect="off"`** -- iOS tries to "correct" emails to real words. Always off.

**`spellCheck={false}`** -- same reason.

---

## 5. Dates & Times

### Date of birth

```tsx
<input
	type="date"
	autoComplete="bday"
	min="1900-01-01"
	max={`${new Date().getFullYear() - 1}-12-31`}
/>
```

- `type="date"` triggers the native date picker -- way better than any custom datepicker on mobile.
- `autoComplete="bday"` (not `birthdate` or `dob`) is the WHATWG standard value.
- `min`/`max` constraints prevent impossible dates (year 0001, future dates for DOB).
- The displayed format follows the device locale, but `value` is always `YYYY-MM-DD`.

### Date with `type="text"` (when you need a custom picker)

Only do this when the native picker doesn't fit (e.g. selecting a date range, booking a slot). Then:

```tsx
<input
  type="text"
  inputMode="numeric"
  placeholder="DD/MM/YYYY"
  pattern="\d{2}/\d{2}/\d{4}"
  value={formatDate(value)}
  onChange={...}
/>
```

### Time

```tsx
<input type="time" autoComplete="off" />
```

Native time picker is solid on both iOS and Android. Don't build a custom one.

### Month

```tsx
<input type="month" />
```

Useful for credit card expiry. Native picker on Android; falls back to text input on iOS Safari (which is acceptable -- `inputMode="numeric"` + `placeholder="MM/YYYY"`).

---

## 6. Address & Geographic

### Standard address fields

```tsx
<input type="text" autoComplete="address-line1" enterKeyHint="next" />
<input type="text" autoComplete="address-line2" enterKeyHint="next" />
<input type="text" autoComplete="address-level2" enterKeyHint="next" /> {/* City */}
<input type="text" autoComplete="address-level1" enterKeyHint="next" /> {/* State/Province */}
<input type="text" inputMode="numeric" autoComplete="postal-code" enterKeyHint="next" />
<select autoComplete="address-country">
  <option value="ID">Indonesia</option>
</select>
```

**Indonesian context**: Province → Regency/City → District → Village (`Provinsi → Kabupaten/Kota → Kecamatan → Kelurahan/Desa`). Use cascading `<select>` elements with autocomplete="off" on each -- the standard `address-level1/2/3` values don't fit the 4-tier Indonesian structure.

### RT/RW

```tsx
<div className="flex gap-2">
	<input
		type="text"
		inputMode="numeric"
		pattern="[0-9]*"
		maxLength={3}
		placeholder="RT"
		autoComplete="off"
		enterKeyHint="next"
	/>
	<input
		type="text"
		inputMode="numeric"
		pattern="[0-9]*"
		maxLength={3}
		placeholder="RW"
		autoComplete="off"
		enterKeyHint="next"
	/>
</div>
```

### Geolocation capture

```tsx
<button onClick={captureLocation}>Ambil Lokasi Saat Ini</button>;

async function captureLocation() {
	if (!navigator.geolocation) {
		alert("Geolokasi tidak didukung perangkat ini");
		return;
	}
	const pos = await new Promise<GeolocationPosition>((resolve, reject) => {
		navigator.geolocation.getCurrentPosition(resolve, reject, {
			enableHighAccuracy: true,
			timeout: 10000,
			maximumAge: 0,
		});
	});
	// pos.coords.latitude, pos.coords.longitude
}
```

- Geolocation requires HTTPS (or localhost).
- On iOS 13+, requires explicit permission prompt -- must be triggered by user gesture.
- Provide a manual-entry fallback for users who deny permission.

---

## 7. File Upload & Camera

### Single image upload

```tsx
<input
	type="file"
	accept="image/*"
	// Don't set capture -- let the user choose between camera and gallery
/>
```

### Force rear camera (KTP scan, etc.)

```tsx
<input
	type="file"
	accept="image/*"
	capture="environment" // "user" for front, "environment" for rear
/>
```

- iOS: opens rear camera directly.
- Android Chrome: opens rear camera directly.
- Some Android launchers ignore `capture` and still show the picker -- that's OK, the constraint is "best effort".

### Multiple files

```tsx
<input type="file" multiple accept="image/*,application/pdf" />
```

### Preview before upload

```tsx
function ImageInput() {
	const [preview, setPreview] = useState<string>();

	return (
		<>
			<input
				type="file"
				accept="image/*"
				onChange={(e) => {
					const file = e.target.files?.[0];
					if (file) setPreview(URL.createObjectURL(file));
				}}
			/>
			{preview && <img src={preview} alt="Pratinjau" className="mt-2 max-h-48 rounded-lg" />}
		</>
	);
}
```

**Cleanup**: `URL.createObjectURL` leaks memory until revoked. Call `URL.revokeObjectURL(preview)` in a `useEffect` cleanup or when replacing the image.

### Camera permission UX

Don't auto-prompt for camera permission on page load -- that's hostile. Let the user tap a "Scan KTP" button first; the browser handles the permission prompt.

---

## 8. OTP / Verification Codes

### Single 6-digit input

```tsx
<input
	type="text"
	inputMode="numeric"
	pattern="[0-9]*"
	autoComplete="one-time-code"
	maxLength={6}
	placeholder="______"
	className="tracking-[0.5em] text-center font-mono text-lg"
/>
```

**`autoComplete="one-time-code"`** is the magic value -- on iOS Safari and Android Chrome, when an SMS with the OTP arrives, the browser suggests it as an autofill option above the keyboard. Users tap once to fill.

### Six separate inputs (one digit per box)

Visually nicer but breaks the autofill. Use this only if you have a strong design reason, and implement:

```tsx
const [digits, setDigits] = useState(["", "", "", "", "", ""]);
const refs = useRef<(HTMLInputElement | null)[]>([]);

const handleChange = (i: number, v: string) => {
	if (!/^\d?$/.test(v)) return;
	const next = [...digits];
	next[i] = v;
	setDigits(next);
	if (v && i < 5) refs.current[i + 1]?.focus();
};

// In JSX
{
	digits.map((d, i) => (
		<input
			key={i}
			ref={(el) => {
				refs.current[i] = el;
			}}
			type="text"
			inputMode="numeric"
			maxLength={1}
			value={d}
			onChange={(e) => handleChange(i, e.target.value)}
			onKeyDown={(e) => {
				if (e.key === "Backspace" && !digits[i] && i > 0) {
					refs.current[i - 1]?.focus();
				}
			}}
			onPaste={(e) => {
				const text = e.clipboardData.getData("text").replace(/\D/g, "").slice(0, 6);
				if (text.length === 6) {
					e.preventDefault();
					setDigits(text.split(""));
					refs.current[5]?.focus();
				}
			}}
			className="w-12 h-14 text-center text-xl font-bold"
		/>
	));
}
```

**Critical**: implement paste -- users copy the OTP from SMS and expect to paste into the first box to fill all six. Without paste support, the UX is worse than a single input.

---

## 9. Multi-Step Forms

### Pattern

```tsx
function MultiStepForm() {
	const [step, setStep] = useState(0);
	const steps = ["Identitas", "Alamat", "Kontak", "Review"];

	return (
		<form>
			{/* Progress indicator */}
			<ol className="flex mb-6">
				{steps.map((label, i) => (
					<li key={i} className="flex-1">
						<div className={`h-1 ${i <= step ? "bg-blue-600" : "bg-slate-200"}`} />
						<span
							className={`text-xs ${i === step ? "text-blue-600 font-medium" : "text-slate-500"}`}
						>
							{label}
						</span>
					</li>
				))}
			</ol>

			{/* Step content */}
			<div className="min-h-[50dvh]">
				{step === 0 && <IdentitasStep />}
				{step === 1 && <AlamatStep />}
				{/* ... */}
			</div>

			{/* Navigation */}
			<div className="flex gap-2 mt-6">
				{step > 0 && (
					<Button variant="outline" onClick={() => setStep(step - 1)} className="flex-1 min-h-11">
						Kembali
					</Button>
				)}
				{step < steps.length - 1 ? (
					<Button onClick={() => setStep(step + 1)} className="flex-1 min-h-11">
						Lanjut
					</Button>
				) : (
					<Button type="submit" className="flex-1 min-h-11">
						Kirim
					</Button>
				)}
			</div>
		</form>
	);
}
```

### WCAG 2.2 SC 3.3.7 -- Redundant Entry compliance

If step 1 asked for `email`, and step 3 needs `email` again, **auto-fill it**:

```tsx
// Step 3
<input
	type="email"
	defaultValue={formData.email} // pre-filled from step 1
	onChange={(e) => updateField("emailConfirmation", e.target.value)}
/>
```

Or offer a dropdown: "Gunakan email yang sama: user@example.com".

### Persistence

Multi-step forms MUST persist data across steps (in `useState` at the parent, or in `localStorage`). If the user navigates away and comes back, the form should restore. The browser's back button should not wipe step 2's data when going back to step 1.

---

## 10. Validation & Error Display

### Inline validation (on blur)

Validate when the user leaves a field, not on every keystroke:

```tsx
function ValidatedInput({ value, onChange, validate }: Props) {
	const [error, setError] = useState<string>();
	const [touched, setTouched] = useState(false);

	return (
		<div>
			<input
				value={value}
				onChange={(e) => {
					onChange(e.target.value);
					if (touched) setError(validate(e.target.value));
				}}
				onBlur={() => {
					setTouched(true);
					setError(validate(value));
				}}
				aria-invalid={!!error}
				aria-describedby={error ? `${id}-error` : undefined}
			/>
			{error && (
				<p id={`${id}-error`} role="alert" className="mt-1 text-sm text-red-600">
					{error}
				</p>
			)}
		</div>
	);
}
```

**`aria-invalid="true"`** tells screen readers the field has an error.
**`aria-describedby`** links the field to its error message -- screen readers announce the error when focusing the field.
**`role="alert"`** on the error paragraph announces it immediately when it appears.

### Submit-time validation

On submit, validate all fields and scroll the first error into view:

```tsx
async function handleSubmit(e: FormEvent) {
	e.preventDefault();
	const errors = validateAll(formData);
	if (Object.keys(errors).length) {
		setErrors(errors);
		// Scroll to first error
		const firstErrorField = document.querySelector("[aria-invalid='true']");
		firstErrorField?.scrollIntoView({ behavior: "smooth", block: "center" });
		firstErrorField?.focus();
		return;
	}
	// submit
}
```

### Don't disable the submit button

A disabled submit button gives no feedback -- users tap it and nothing happens, so they think the app is broken. Instead, always allow the tap, then show errors. The exception is during async submit -- disable to prevent double-submission, and show a loading spinner on the button.
