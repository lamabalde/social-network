import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import AllPosts from "./AllPosts";
// jest.setTimeout(10000);
describe("Posts and Comments feature", () => {
	test("Loading Post", () => {
		render(
			<MemoryRouter>
				<AllPosts isLoadingPost={true} />
			</MemoryRouter>
		);

		// Act
		// Assert
		const loadingMsg = screen.getByText(/Loading/i, {
			exact: false,
		});
		expect(loadingMsg).toBeInTheDocument();
	});

	test("Get All Posts (No Posts)", async () => {
		// Arrange
		window.fetch = jest.fn();
		window.fetch.mockResolvedValueOnce({
			json: async () => [],
		});

		render(
			<MemoryRouter>
				<AllPosts />
			</MemoryRouter>
		);

		// Act
		// Assert
		const postSomethingMsg = await screen.findByText(
			/No Posts Yet/i,
			{
				exact: false,
			},
			{
				timeout: 5000,
			}
		);
		expect(postSomethingMsg).toBeInTheDocument();
	});

	test("Get All Posts (>=1 Posts)", async () => {
		// Arrange
		window.fetch = jest.fn();
		window.fetch.mockResolvedValueOnce({
			json: async () => {
				return [
					{
						author: 1,
						avatar: "",
						createdat: "2023-08-01 03:16:00.000000 +0100 +0100",
						fname: "testfname1",
						id: 1,
						image: "",
						lname: "testlname1",
						message: "post1_testmsg",
						nname: "testnname1",
						privacy: 1,
						success: true,
					},
					{
						author: 1,
						avatar: "",
						createdat: "2023-08-01 15:16:00.000000 +0100 +0100",
						fname: "testfname1",
						id: 2,
						image: "",
						lname: "testlname1",
						message: "post2_testmsg",
						nname: "testnname1",
						privacy: 0,
						success: true,
					},
				];
			},
		});

		render(
			<MemoryRouter>
				<AllPosts />
			</MemoryRouter>
		);

		// Act

		// Assert
		const fname = await screen.findByText(
			"testfname1",
			{ exact: true },
			{ timeout: 3000 }
		);
		screen.debug(fname);
		expect(fname).toBeInTheDocument();
		const lname = await screen.findByText(
			"testlname1",
			{ exact: true },
			{ timeout: 3000 }
		);
		expect(lname).toBeInTheDocument();
		const nname = await screen.findByText(
			"testnname1",
			{ exact: true },
			{ timeout: 3000 }
		);
		expect(nname).toBeInTheDocument();
		//9 Aug 23, 18:36
		const post1createdat = await screen.findByText(
			"1 Aug 23, 03:16",
			{
				exact: true,
			},
			{ timeout: 3000 }
		);
		expect(post1createdat).toBeInTheDocument();
		const post2createdat = await screen.findByText(
			"1 Aug 23, 15:16",
			{
				exact: true,
			},
			{ timeout: 3000 }
		);
		expect(post2createdat).toBeInTheDocument();
		const post1msg = await screen.findByText(
			"post1_testmsg",
			{ exact: true },
			{ timeout: 3000 }
		);
		expect(post1msg).toBeInTheDocument();
		const post2msg = await screen.findByText(
			"post2_testmsg",
			{ exact: true },
			{ timeout: 3000 }
		);
		expect(post2msg).toBeInTheDocument();
	});
});
