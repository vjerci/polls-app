import { useAppSelector } from "@/store/hooks";
import { selectUserName } from "@/store/reducers/user";
import LogoutButton from "@/containers/polls/logout_button";

export default function ListHeader({
	title,
	userName,
}: {
	title: string;
	userName: string;
}) {
	return (
		<div className="pb-4">
			<nav className="mb-4 grid">
				<LogoutButton />
			</nav>

			<span className="text-xl text-stone-100 pb-4 block">
				Welcome back {userName}
			</span>

			<h2 className="text-2xl text-slate-100">{title}</h2>
		</div>
	);
}
