// MIT License
// Copyright (c) 2020 Luis Espino

function reflex_agent(location, state) {
    if (state == "DIRTY") return "CLEAN";
    else if (location == "A") return "RIGHT";
    else if (location == "B") return "LEFT";
}

function test(states) {
    //ubica donde esta posicionado
    var location = states[0];

    /*   Imprime el estado actual    */
    // el * representa la posicion de la aspiradora
    var locA = location == "A" ? "| * " : "|   ";
    var locB = location == "A" ? " |" : " * |";

    //Mensaje de salida para mostrar el estado actual.
    document.getElementById("log").innerHTML += "<br>[ ".concat(locA).concat(states[1]).concat(" | ").concat(states[2]).concat(locB).concat(" ]");

    //condicion ? valor si verdadero : valor si falso
    var state = states[0] == "A" ? states[1] : states[2];
    //ejecuta la funcion que evalua el estado y decide si limpia o mueve a izquierda o derecha
    var action_result = reflex_agent(location, state);
    //imprime la accion
    //document.getElementById("log").innerHTML += "<br>Location: ".concat(location).concat(" | Action: ").concat(action_result);
    if (action_result == "CLEAN") {
        if (location == "A"){
            states[1] = "CLEAN";
            countstates++;
        } 
        else if (location == "B"){
            states[2] = "CLEAN";
            countstates++;
        }
    }
    else {
        if (states[0] == "A" && states[2] == "CLEAN" && states[1] == "CLEAN"){
            states[0] = "B";
            states[1] = "DIRTY";
            states[2] = "DIRTY";
            countstates++;
        }
        else if (action_result == "RIGHT"){
            states[0] = "B";
            countstates++;
        } 
        else if (action_result == "LEFT"){
            states[0] = "A";
            countstates++;
        } 
    }
    if (countstates < 9) setTimeout(function () { test(states); }, 2000);
}

var countstates = 0;
var states = ["A","DIRTY","DIRTY"];
test(states);

//Muestra de como son los otrso states posibles
//var states = ["A","DIRTY","CLEAN"];
//var states = ["A","CLEAN","DIRTY"];
//var states = ["A","CLEAN","CLEAN"];
//var states = ["B","DIRTY","DIRTY"];
//var states = ["B","DIRTY","CLEAN"];
//var states = ["B","CLEAN","DIRTY"];
//var states = ["B","CLEAN","CLEAN"];

