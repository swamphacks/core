/**
 *
 * @param values - form values
 * @param fieldMetaMap - map of field metadata
 * @returns - only the dirty fields from values
 */
export function pickDirty<T extends object>(
  values: T,
  fieldMetaMap: Record<keyof T, { isDirty: boolean }>,
): Partial<T> {
  const result = {} as Partial<T>;
  for (const key in values) {
    if (Object.prototype.hasOwnProperty.call(values, key)) {
      const meta = fieldMetaMap[key as keyof T];
      const val = values[key as keyof T];
      if (meta?.isDirty) {
        result[key as keyof T] = val;
      }
    }
  }
  return result;
}
