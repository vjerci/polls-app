import Error from "@/components/common/error";

export default function Input({
	placeholder,
	onChange,
	value,
	error,
}: {
	placeholder: string;
	onChange: any;
	value: string;
	error?: string;
}) {
	return (
		<div>
			<label
				htmlFor={placeholder}
				className="text-stone-100 cursor-pointer"
			>
				{placeholder} :
			</label>
			<input
				id={placeholder}
				placeholder={placeholder}
				onChange={onChange}
				value={value}
				className="mt-2 mb-4 px-2 py-1 rounded w-full text-stone-900"
			/>
			<Error error={error} />
		</div>
	);
}
