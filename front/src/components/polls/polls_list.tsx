import { Poll } from "@/store/reducers/polls_list";
import Link from "@/components/common/link";
import { capitalize } from "@/lib/text";
import Loader from "../common/loader";

export default function PollsList({
	polls,
	loading,
	pollClicked,
}: {
	polls: Array<Poll>;
	loading: boolean;
	pollClicked: (id: string) => any;
}): React.JSX.Element {
	if (polls.length == 0) {
		return (
			<h2 className="text-xl text-stone-500 py-8">
				{" "}
				There are no Polls created yet
			</h2>
		);
	}

	if (loading) {
		return (
			<div className="py-8 relative h-[428px]">
				<div className="absolute top-0 right-0 left-0 bottom-0 width-100 height-100 grid">
					<div className="place-self-center">
						<Loader />
					</div>
				</div>
			</div>
		);
	}

	return (
		<div className="py-8 h-[428px]">
			{polls.map((p) => (
				<Link
					key={p.id}
					className="block py-1 px-2 my-1 bg-slate-800 hover:bg-slate-900"
					onClick={() => {
						pollClicked(p.id);
					}}
				>
					{capitalize(p.name)}
				</Link>
			))}
		</div>
	);
}
