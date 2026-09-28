import { myHp, print } from "kolmafia";
import { $item } from "libram";

// This file is a compile-time smoke target showing the imports generated ASH
// is expected to use. It is not an autonomous executor.
export function runtimeSmoke(): void {
  const hp = myHp();
  const tooth = $item`seal tooth`;
  print(`kingdomsitter TS target loaded; hp=${hp}; sample=${tooth}`);
}
