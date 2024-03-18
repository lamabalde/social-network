import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import RegForm from "./RegForm";

describe("Reg Form Compo", () => {
    test("renders Email Input", () => {
        // Arrange
        render(
            <MemoryRouter>
                <RegForm />
            </MemoryRouter>
        );
        // Act
        // Assert
        const emailInput = screen.getByLabelText("Email");
        expect(emailInput).toBeInTheDocument();
    });
    
    test("renders Password Input", () => {
        // Arrange
        render(
            <MemoryRouter>
                <RegForm />
            </MemoryRouter>
        );
        // Act
        // Assert
        const pwInput = screen.getByLabelText("Password", {exact: true});
        expect(pwInput).toBeInTheDocument();
    });

    // test("");

    // test("");

    // test("");

    // test("");

    // test("");

    // test("");

    test("render reg btn", () => {
        // Arrange
        render(
            <MemoryRouter>
                <RegForm />
            </MemoryRouter>
        );
        // Act
        //Assert
        const regBtn = screen.getByRole("button", { name: /register/i});
        expect(regBtn).toBeInTheDocument();
    });

    test("renders login link", () => {
        // Arrange
        render(
            <MemoryRouter>
                <RegForm />
            </MemoryRouter>
        );
        // Act
        // Assert
        const loginLinkEl = screen.getByRole("link", { name: /Login/i });
        expect(loginLinkEl).toBeInTheDocument();
    });
});