import { Poll } from "@/store/reducers/polls_list";
import Link from "@/components/common/link";
import { capitalize } from "@/lib/text";
import next from "next";

export default function ListPagination({
	hasNext,
	page,
}: {
	hasNext: boolean;
	page: number;
}): React.JSX.Element {
	if (!hasNext && page == 0) {
		return (
			<h2 className="text text-stone-500 py-4">
				There is no next page, create more polls to enable pagination
			</h2>
		);
	}

	const prevPage = page - 1;
	const nextPage = page + 1;

	return (
		<div className="py-4 columns-3">
			{page != 0 ? (
				<Link href={`/polls?page=${prevPage}`} className="inline-block">
					Previous page
				</Link>
			) : (
				<div>&nbsp;</div>
			)}

			<div className="text-stone-500 text-center">{page}</div>

			{hasNext ? (
				<Link
					href={`/polls?page=${nextPage}`}
					className="inline-block float-right"
				>
					Next Page
				</Link>
			) : (
				<></>
			)}
		</div>
	);
}
