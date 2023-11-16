import Breadcrumb from "@/components/common/breadcrumb";
import { selectPollListPage } from "@/store/reducers/polls_list";
import { useSelector } from "react-redux";

export default function PollBreadCrumb({
	currentName,
}: {
	currentName: string;
}) {
	const page = useSelector(selectPollListPage);

	return (
		<Breadcrumb
			homeHref={`/polls?page=${page}`}
			links={[{ name: currentName }]}
		></Breadcrumb>
	);
}
