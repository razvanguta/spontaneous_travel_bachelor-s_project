
fetch("/jsonReview/" + document.getElementById("theAgencyID").innerHTML)
.then(response => response.json())
.then(data => {
    var nr = 0
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        const divAll = document.createElement("div");
        divAll.setAttribute("id","div"+i);
        divAll.setAttribute("class","formular-login2");
        document.getElementById("listOfReviews").appendChild(divAll);
        var child = document.getElementById("div"+i);
        //add the text
        const paragraph = document.createElement("h2");
        paragraph.setAttribute("id","client-paragraph"+i);
        child.appendChild(paragraph);
        paragraph.innerHTML="Clientul "+obj.client +" a oferit urmatoarea recenzie:";

        //add the text
        const paragraph2 = document.createElement("h3");
        paragraph2.setAttribute("id","title-paragraph"+i);
        child.appendChild(paragraph2);
       paragraph2.innerHTML=obj.title;

        //add the text
        const paragraph3 = document.createElement("p");
        paragraph3.setAttribute("id","comm-paragraph"+i);
        child.appendChild(paragraph3);
        document.getElementsByTagName("p")[nr].innerHTML=obj.comment;
        nr = nr + 1;
            //add the text
          for(var j=0;j< parseInt(obj.stars); j++){
              console.log(1);
              const star = document.createElement("SPAN");
              star.setAttribute("class","s"+j)
              star.innerHTML="☆";
              child.appendChild(star);
          }
        
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
        var form = document.createElement("form");
        form.setAttribute("id","deleteReview");
        form.method="get";
        form.action="/deleteReview/"+ document.getElementById("theUserID").innerHTML + "/" + document.getElementById("theAgencyID").innerHTML+"/"+obj.date;  
        var button = document.createElement("button");
        button.innerHTML = "Sterge recenzia";
        button.setAttribute("id","delete-button"+i);
        button.setAttribute("class","bbtnn");
        button.setAttribute("class","toRegister-btn");
        form.appendChild(button);
        child.append(form);

        // var button2 = document.createElement("a");
        // button2.innerHTML = "Editeaza recenzia";
        // button2.setAttribute("id","update-button"+i);
        // button2.href = "/updateReview/"+ document.getElementById("theUserID").innerHTML + "/" + document.getElementById("theAgencyID").innerHTML+"/"+obj.date;   
        // child.appendChild(button2); 
        }
        child.appendChild(separator); 
    }
}) 