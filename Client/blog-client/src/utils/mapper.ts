export const getDirtyFieldNames = (dirtyFields: object, prefix = "") => {
  return Object.entries(dirtyFields).flatMap(([key, value]): unknown => {
    const name = prefix ? `${prefix}.${key}` : key;

    if (value === true) {
      return name;
    }

    if (typeof value === "object") {
      return getDirtyFieldNames(value, name);
    }

    return [];
  });
};

export function getTypeValidValue<const T extends readonly string[]>(
  val: string | null,
  validVals: T,
  defaultVal: T[number]
): T[number] {
  if (val && validVals.includes(val as T[number])) {
    return val as T[number];
  }
  return defaultVal;
}

export function getQueryParam(object: unknown) {
  const params = new URLSearchParams();
  if (!object) return params;
  Object.entries(object).forEach(([key, value]) => {
    if (value === undefined || value === null) return;

    if (Array.isArray(value)) {
      value.forEach((v) => params.append(key, v));
    } else {
      params.append(key, value.toString());
    }
  });
  return params;
}
