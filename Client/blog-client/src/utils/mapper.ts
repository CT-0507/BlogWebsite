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
