import Loader from "@/components/common/loader";
import DetailsAnswer from "@/components/polls/details_answer";
import LogoutButton from "@/containers/polls/logout_button";
import PollBreadCrumb from "@/containers/polls/poll_breadcrumb";
import { capitalize } from "@/lib/text";
import { useAppDispatch } from "@/store/hooks";
import {
	fetchPollDetails,
	pollDetailsVote,
	selectPollDetailsAnswers,
	selectPollDetailsIsLoading,
	selectPollDetailsName,
} from "@/store/reducers/poll_details";
import { useRouter } from "next/router";
import { useEffect } from "react";
import { useSelector } from "react-redux";

export default function Page() {
	const dispatch = useAppDispatch();
	const router = useRouter();

	useEffect(() => {
		const idS = router.query.id as string;
		if (idS) {
			dispatch(fetchPollDetails({ pollId: idS }));
		}
	}, [dispatch, router.query.id]);

	const loading = useSelector(selectPollDetailsIsLoading);
	const name = useSelector(selectPollDetailsName);
	const answers = useSelector(selectPollDetailsAnswers);

	const vote = (id: string) => {
		dispatch(pollDetailsVote({ answerId: id }));
	};

	return (
		<>
			<nav className="mb-4 grid">
				<LogoutButton />
			</nav>
			<PollBreadCrumb currentName="Poll details" />
			<h2 className="text-2xl text-slate-100 pb-4">{capitalize(name)}</h2>

			{loading ? (
				<div className="py-8 relative h-[200px]">
					<div className="absolute top-0 right-0 left-0 bottom-0 width-100 height-100 grid">
						<div className="place-self-center">
							<Loader />
						</div>
					</div>
				</div>
			) : (
				<>
					{answers?.map((answer) => (
						<DetailsAnswer
							key={answer.id}
							answer={answer.name}
							id={answer.id}
							hasVoted={answer.hasUserVoted}
							votesCount={answer.votesCount}
							onClick={vote}
						/>
					))}
				</>
			)}
		</>
	);
}
