export default function Error({ error }: { error?: string }) {
	if (error !== "") {
		return <span className="mt-2 mb-4 text-red-500">{error}</span>;
	}

	return <></>;
}
