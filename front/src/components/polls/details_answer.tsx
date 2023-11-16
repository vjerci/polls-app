import { capitalize } from "@/lib/text";
import Link from "../common/link";

export default function DetailsAnswer({
	id,
	answer,
	votesCount,
	hasVoted,
	onClick,
}: {
	id: string;
	answer: string;
	votesCount: number;
	hasVoted?: boolean;
	onClick: (id: string) => any;
}) {
	let disabledClassNames = "";
	if (hasVoted) {
		disabledClassNames =
			"bg-slate-700 hover:text-link hover:bg-slate-700 cursor-default";
	}

	const clicked = () => {
		if (hasVoted) {
			return;
		}

		onClick(id);
	};

	return (
		<Link
			className={`flex items-baseline py-1 px-2 my-1 bg-slate-800 hover:bg-slate-900 ${disabledClassNames}`}
			onClick={clicked}
		>
			<span className="block">{capitalize(answer)}</span>
			{hasVoted ? <span className="ml-1 text-sm">(voted)</span> : <></>}

			<span className="ml-auto block text-slate-500 text-sm">
				{" "}
				({votesCount == 1 ? `1 vote` : `${votesCount} votes`}){" "}
			</span>
		</Link>
	);
}
