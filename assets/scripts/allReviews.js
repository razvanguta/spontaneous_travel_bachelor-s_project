
fetch("/jsonReview/" + document.getElementById("theAgencyID").innerHTML)
.then(response => response.json())
.then(data => {
    var nr = 0
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        const divAll = document.createElement("div");
        divAll.setAttribute("id","div"+i);
        document.getElementById("listOfReviews").appendChild(divAll);
        var child = document.getElementById("div"+i);
        //add the text
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","client-paragraph"+i);
        child.appendChild(paragraph);
        document.getElementsByTagName("p")[nr].innerHTML="Clientul "+obj.client +" a oferit urmatoarea recenzie:";
        nr = nr + 1;
        //add the text
        const paragraph2 = document.createElement("p");
        paragraph2.setAttribute("id","title-paragraph"+i);
        child.appendChild(paragraph2);
        document.getElementsByTagName("p")[nr].innerHTML="Titlul: " + obj.title;
        nr = nr + 1;
        //add the text
        const paragraph3 = document.createElement("p");
        paragraph3.setAttribute("id","comm-paragraph"+i);
        child.appendChild(paragraph3);
        document.getElementsByTagName("p")[nr].innerHTML="Comentariul: " + obj.comment;
        nr = nr + 1;
            //add the text
        const paragraph7 = document.createElement("p");
        paragraph7.setAttribute("id","stars-paragraph"+i);
        child.appendChild(paragraph7);
        document.getElementsByTagName("p")[nr].innerHTML="Nota oferita: "+ obj.stars;
        nr = nr + 1
        //add the text
        const paragraph6 = document.createElement("p");
        paragraph6.setAttribute("id","date-paragraph"+i);
        child.appendChild(paragraph6);
        document.getElementsByTagName("p")[nr].innerHTML="Recenize lasata la data de " + obj.date;
        nr = nr + 1;
        //add space
        var space = document.createElement("br");
        var separator = document.createElement("hr");
        child.appendChild(space);
        //add a button
        if(obj.same == "yes"){
        var button = document.createElement("a");
        button.innerHTML = "Sterge recenzia";
        button.setAttribute("id","delete-button"+i);
        button.href = "/deleteReview/"+ document.getElementById("theUserID").innerHTML + "/" + document.getElementById("theAgencyID").innerHTML+"/"+obj.date;  
        child.appendChild(button);

        child.append(space)

        // var button2 = document.createElement("a");
        // button2.innerHTML = "Editeaza recenzia";
        // button2.setAttribute("id","update-button"+i);
        // button2.href = "/updateReview/"+ document.getElementById("theUserID").innerHTML + "/" + document.getElementById("theAgencyID").innerHTML+"/"+obj.date;   
        // child.appendChild(button2); 
        }
        child.appendChild(separator); 
    }
}) 