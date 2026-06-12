function truncate(text: string, max = 100) {
  return text.length > max ? text.slice(0, max) + "..." : text;
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString();
}

function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export { truncate, formatDate, capitalize };
