import "./App.css";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Header from "./components/Header";
import SignIn from "./templates/SignIn";
import SignUp from "./templates/SignUp";
import CustomizedInputBase from "./components/Search";
import ApiFecth from "./components/ApiFetch";
import ChartLists from "./components/ChartLists";
import Dashboard from "./components/Dashboard";
import { Shops } from "./components/Shop";
import { Items } from "./components/Item";

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Header />
        <Routes>
          <Route path="/sign-in" element={<SignIn />} />
          <Route path="/sign-up" element={<SignUp />} />
          <Route path="/list" element={<ApiFecth />} />
          <Route path="/chartList/*" element={<ChartLists />} />
          <Route path="/" element={<CustomizedInputBase />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/shops" element={<Shops />} />
          <Route path="/items" element={<Items />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
