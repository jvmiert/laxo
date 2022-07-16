export const formatPrice = (value: number, name: string): string => {
  return value ? value.toLocaleString("vi-VN") : "";
};

export const parsePrice = (value: string, name: string): number => {
  return parseFloat(value.replaceAll(".", "")) || 0;
};
