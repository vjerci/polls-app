import Link from "@/components/common/link";
import { type } from "os";

export default function Breadcrumb({
	homeHref,
	links,
}: {
	homeHref: string;
	links: Array<{ href?: string; name: string }>;
}) {
	return (
		<nav className="my-4">
			<Link href={homeHref}> Home </Link>

			{links.map((link, index) => (
				<span key={index}>
					<span className="mx-4">-&gt;</span>
					{typeof link.href === "undefined" ? (
						<span className="text-stone-100">{link.name}</span>
					) : (
						<Link href={link.href}> {link.name} </Link>
					)}
				</span>
			))}
		</nav>
	);
}
