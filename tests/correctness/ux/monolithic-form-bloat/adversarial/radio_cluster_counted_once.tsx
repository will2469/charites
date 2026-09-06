import React from 'react';

export function RadioClusterForm() {
  return (
    <form className="space-y-4">
      {/* 8 radio inputs sharing name="tier" count as 1 logical field + 2 text inputs = 3 fields <= 9 */}
      <div className="space-y-1">
        <label><input type="radio" name="tier" value="1" /> Tier 1</label>
        <label><input type="radio" name="tier" value="2" /> Tier 2</label>
        <label><input type="radio" name="tier" value="3" /> Tier 3</label>
        <label><input type="radio" name="tier" value="4" /> Tier 4</label>
        <label><input type="radio" name="tier" value="5" /> Tier 5</label>
        <label><input type="radio" name="tier" value="6" /> Tier 6</label>
        <label><input type="radio" name="tier" value="7" /> Tier 7</label>
        <label><input type="radio" name="tier" value="8" /> Tier 8</label>
      </div>

      <input name="contact_person" placeholder="Contact Person" />
      <input name="email_address" placeholder="Email Address" />
      <button type="submit">Submit</button>
    </form>
  );
}
