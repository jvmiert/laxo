export function generatePaginateNumbers(
  current: number,
  total: number,
): Array<string | number> {
  const center: Array<string | number> = [
    current - 2,
    current - 1,
    current,
    current + 1,
    current + 2,
  ];
  const filteredCenter = center.filter((p) => p > 1 && p < total);
  const includeThreeLeft = current === 5;
  const includeThreeRight = current === total - 4;
  const includeLeftDots = current > 5;
  const includeRightDots = current < total - 4;

  if (includeThreeLeft) filteredCenter.unshift(2);
  if (includeThreeRight) filteredCenter.push(total - 1);

  if (includeLeftDots) filteredCenter.unshift("...");
  if (includeRightDots) filteredCenter.push("...");

  if (total == 1) {
    return [1];
  }
  return [1, ...filteredCenter, total];
}
