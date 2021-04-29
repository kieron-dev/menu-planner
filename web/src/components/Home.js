import React from "react";
import { useAuth } from "./AuthProvider";
import { Menu, Container, Image, Icon, Dropdown } from "semantic-ui-react";

const Home = () => {
    const { auth, logout } = useAuth();

    return (
        <>
            <Menu inverted stackable>
                <Container>
                    <Menu.Item>
                        <Image
                            size="mini"
                            src="/logo192.png"
                            style={{ marginRight: "1.5em" }}
                        />
                        Menu Planner
                    </Menu.Item>
                    <Menu.Item as="a">Home</Menu.Item>
                    <Menu.Item as="a">Planning</Menu.Item>
                    <Menu.Item as="a">Recipes</Menu.Item>
                    <Menu.Item as="a">Shopping List</Menu.Item>
                    <Menu.Menu position="right">
                        <Dropdown text={auth.name} item>
                            <Dropdown.Menu>
                                <Dropdown.Item as="a">
                                    <Icon name="user" />
                                    Profile
                                </Dropdown.Item>
                                <Dropdown.Item as="a">
                                    <Icon name="setting" />
                                    Settings
                                </Dropdown.Item>
                                <Dropdown.Divider />
                                <Dropdown.Item as="a" onClick={logout}>
                                    <Icon name="log out" />
                                    Logout
                                </Dropdown.Item>
                            </Dropdown.Menu>
                        </Dropdown>
                    </Menu.Menu>
                </Container>
            </Menu>
        </>
    );
};

export default Home;
