export const capitalize = (str: string): string => {
	if (str.length > 1) {
		return str[0].toUpperCase() + str.slice(1);
	}

	return str;
};
