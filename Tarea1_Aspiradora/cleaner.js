// MIT License
// Copyright (c) 2020 Luis Espino

function reflex_agent(location, state) {
    if (state == "DIRTY") return "CLEAN";
    else if (location == "A") return "RIGHT";
    else if (location == "B") return "LEFT";
}

function test(states, visitedStates) {
    var location = states[0];
    var state = states[0] == "A" ? states[1] : states[2];
    var action_result = reflex_agent(location, state);

    // Log the action taken
    document.getElementById("log").innerHTML += "<br>Location: " + location + " | Action: " + action_result;

    // Mark the state as CLEAN if action is CLEAN
    if (action_result == "CLEAN") {
        if (location == "A") states[1] = "CLEAN";
        else if (location == "B") states[2] = "CLEAN";
    }
    // Move to the next location based on action
    else if (action_result == "RIGHT" && states[0] == "A") states[0] = "B";
    else if (action_result == "LEFT" && states[0] == "B") states[0] = "A";		

    // Record the visited state (a combination of location and cleanliness of both places)
    let currentState = states[0] + states[1] + states[2];
    if (!visitedStates.includes(currentState)) {
        visitedStates.push(currentState); // only add unique states
    }

    // Check if we have visited all 8 states (A clean, B clean in all combinations)
    if (visitedStates.length >= 8) {
        document.getElementById("log").innerHTML += "<br>All states visited. Stopping.";
        return;
    }

    // Continue the test after 2 seconds
    setTimeout(function () { test(states, visitedStates); }, 2000);
}

// Initial states and visited states tracking
var states = ["A", "DIRTY", "DIRTY"];
var visitedStates = [];
test(states, visitedStates);
