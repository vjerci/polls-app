import { FormEventHandler } from "react";
import Loader from "./loader";

export default function Form({
	children,
	submitName,
	onSubmit,
	loading,
}: {
	children: Array<React.ReactNode> | React.ReactNode;
	submitName: string;
	onSubmit: FormEventHandler<HTMLFormElement>;
	loading: boolean;
}) {
	return (
		<form
			className="mt-6 bg-slate-800 rounded-xl shadow-lg p-4 w-auto sm:w-[400px] grid border border-slate-700 relative"
			onSubmit={onSubmit}
		>
			{children}
			<button
				type="submit"
				className="mt-6 bg-indigo-500 text-stone-100 rounded px-2 py-1 hover:bg-indigo-600 hover:text-stone-200 justify-self-end"
			>
				{submitName}
			</button>

			{loading ? (
				<div className="absolute top-0 right-0 left-0 bottom-0 width-100 height-100 grid bg-opacity-60 bg-slate-900">
					<div className="place-self-center">
						<Loader />
					</div>
				</div>
			) : (
				<></>
			)}
		</form>
	);
}
