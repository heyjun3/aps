import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import { Button } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import MenuIcon from "@mui/icons-material/Menu";
import { Link, useNavigate } from "react-router-dom";

export default function ButtonAppBar() {
  const navigate = useNavigate();
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}
          >
            <MenuIcon />
          </IconButton>
          <Typography component="div" sx={{ flexGrow: 1 }}>
            <Button onClick={() => navigate("/")} color="inherit">
              home
            </Button>
            <Button onClick={() => navigate("/list")} color="inherit">
              list
            </Button>
            <Button onClick={() => navigate("/items")} color="inherit">
              item
            </Button>
            <Button onClick={() => navigate("/shops")} color="inherit">
              shop
            </Button>
          </Typography>
          <Button component={Link} to="/sign-in" color="inherit">
            Sign in
          </Button>
        </Toolbar>
      </AppBar>
    </Box>
  );
}
