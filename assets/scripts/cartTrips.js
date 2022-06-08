
fetch("/jsonSeeCart/")
.then(response => response.json())
.then(data => {
    var nr = 0
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        const divAll = document.createElement("div");
        divAll.setAttribute("id","div"+i);
        divAll.setAttribute("class","formular-login2")
        document.getElementById("cartTrips").appendChild(divAll);
        var child = document.getElementById("div"+i);
        //add the text
        const paragraph = document.createElement("h2");
        paragraph.setAttribute("id","title-paragraph"+i);
        child.appendChild(paragraph);
        paragraph.innerHTML=obj.title;
        nr = nr + 1;
        //add a button
        var form = document.createElement("form");
        form.setAttribute("id","buy");
        form.method="get";
        form.action="/buyTripPage/"+ obj.tripId;  
        var button = document.createElement("button");
        button.innerHTML = "Cumpara excursia";
        button.setAttribute("id","buy-button"+i);
        button.setAttribute("class","toRegister-btn")
        button.href = "/buyTripPage/"+ obj.tripId;  
        form.appendChild(button);
        child.appendChild(form);
         //add a button
         var form2 = document.createElement("form");
        form2.setAttribute("id","out");
        form2.method="get";
        form2.action="/outFromCart/"+ obj.id;   
         var button2 = document.createElement("button");
         button2.innerHTML = "Scoate din cosul de cumparaturi";
         button2.setAttribute("id","out-button"+i);
         button2.setAttribute("class","toRegister-btn")
         form2.appendChild(button2);
         child.appendChild(form2);
        var separator = document.createElement("hr");
        child.appendChild(separator); 
    }
}) 