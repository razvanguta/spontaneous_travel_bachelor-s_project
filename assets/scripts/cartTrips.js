
fetch("/jsonSeeCart/")
.then(response => response.json())
.then(data => {
    var nr = 0
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        const divAll = document.createElement("div");
        divAll.setAttribute("id","div"+i);
        document.getElementById("cartTrips").appendChild(divAll);
        var child = document.getElementById("div"+i);
        //add the text
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","title-paragraph"+i);
        child.appendChild(paragraph);
        document.getElementsByTagName("p")[nr].innerHTML="Titlul: " + obj.title;
        nr = nr + 1;
        //add a button
        var button = document.createElement("a");
        button.innerHTML = "Cumpara excursia";
        button.setAttribute("id","buy-button"+i);
       // button.href = "/deleteReview/"+ document.getElementById("theUserID").innerHTML + "/" + document.getElementById("theAgencyID").innerHTML+"/"+obj.date;  
        child.appendChild(button);
         //add a button
         var button2 = document.createElement("a");
         button2.innerHTML = "Scoate din cosul de cumparaturi";
         button2.setAttribute("id","out-button"+i);
         button2.href = "/outFromCart/"+ obj.id;  
         child.appendChild(button2);
        var separator = document.createElement("hr");
        child.appendChild(separator); 
    }
}) 