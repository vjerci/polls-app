import Link from "@/components/common/link";
import ListHeader from "@/components/polls/list_header";
import ListPagination from "@/components/polls/list_pagination";
import PollsList from "@/components/polls/polls_list";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { addBasicPollInfo } from "@/store/reducers/poll_details";
import {
	fetchPollsList,
	selectPollsListData,
	selectPollsListHasNextPage,
	selectPollsListIsLoading,
} from "@/store/reducers/polls_list";
import { selectUserName } from "@/store/reducers/user";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/router";
import { useEffect } from "react";

export default function Page() {
	const router = useRouter();
	const dispatch = useAppDispatch();
	const searchParams = useSearchParams();

	const pageStr = searchParams.get("page");
	const page = pageStr ? parseInt(pageStr, 10) : 0;

	useEffect(() => {
		dispatch(fetchPollsList({ page: page }));
	}, [dispatch, page]);

	const polls = useAppSelector(selectPollsListData);
	const hasNextPage = useAppSelector(selectPollsListHasNextPage);

	const userName = useAppSelector(selectUserName);
	const isLoading = useAppSelector(selectPollsListIsLoading);

	const onPollClicked = (pollId: string) => {
		const poll = polls.find((p) => p.id === pollId);
		if (poll) {
			dispatch(addBasicPollInfo({ name: poll.name }));
		}
		router.push(`/polls/${pollId}`);
	};

	const refreshData = () => {
		dispatch(fetchPollsList({ page: page }));
	};

	return (
		<>
			<ListHeader title="Latest Pools" userName={userName} />
			<Link onClick={refreshData} className="mr-2">
				Refresh data
			</Link>
			<span className="text-slate-500">|</span>
			<Link href="/polls/create" className="ml-2">
				Create poll
			</Link>
			<PollsList
				polls={polls}
				loading={isLoading}
				pollClicked={onPollClicked}
			/>
			<ListPagination hasNext={hasNextPage} page={page} />
		</>
	);
}
