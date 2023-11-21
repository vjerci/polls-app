export default function Button({
	name,
	onClick,
	className,
}: {
	name: string;
	onClick?: any;
	className?: string;
}) {
	return (
		<button
			className={`mb-5 bg-indigo-500 text-stone-100 rounded px-2 py-1 mt-2 hover:bg-indigo-600 hover:text-stone-200 ${className}`}
			onClick={(e) => {
				e.preventDefault();
				if (onClick) {
					onClick();
				}
			}}
		>
			{name}
		</button>
	);
}
