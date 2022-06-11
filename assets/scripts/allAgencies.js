fetch("/jsonAllAgencies/")
.then(response => response.json())
.then(data => {
    var body = document.getElementById("body");
    
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];

        //add the text
        const paragraph = document.createElement("h2");
        paragraph.setAttribute("id","agency-paragraph");
        paragraph.setAttribute("class","card-text");
 
        paragraph.innerHTML="Agentia de turism "+obj.username;
        var divCenter = document.createElement("div");
        divCenter.setAttribute("class","center2")
        //add the photo
        var img = document.createElement("img");
        img.setAttribute("id","agency-photo");  
        img.setAttribute("class","card-img-top")    
    
        img.src = obj.profile_image
        divCenter.appendChild(img);
        
        //add button
        var form = document.createElement("form");
        form.setAttribute("id","agencyPage");
        form.method="get";
        form.action="/agencyPage/"+obj.username;
        var button = document.createElement("button");
        button.innerHTML = "Viziteaza pagina";
        button.type = "submit";
        button.setAttribute("id","agency-button");
        button.setAttribute("class","toRegister-btn")
        form.appendChild(button);

        //add space
        var space = document.createElement("br");
        var form2 = document.createElement("form");
        form2.setAttribute("id","seeReviews");
        form2.method="get";
        form2.action="/seeReviews/"+obj.username;
        body.appendChild(space);
        var button2 = document.createElement("button");
        button2.innerHTML = "Vizualizeaza sau adauga recenzii";
        button2.setAttribute("id","review-button");
        button2.setAttribute("class","toRegister-btn")
        form2.appendChild(button2);


        if(obj.is_admin == "yes"){
        var form3 = document.createElement("form");
        form3.setAttribute("id","deleteAgency");
        form3.method="get";
        form3.action="/deleteAgency/"+obj.id;
        var button3 = document.createElement("button");
        button3.innerHTML = "Sterge agentia";
        button3.setAttribute("id","review-button");
        button3.setAttribute("class","toRegister-btn")
        form3.appendChild(button3);

        }
        
        var inside = document.createElement("div");
        inside.setAttribute("class","formular-login");
       // inside.setAttribute("style","width: 15%");
        var inside2 = document.createElement("div");
        inside.setAttribute("class","formular-login2");
        
        inside.appendChild(paragraph);
        inside.appendChild(divCenter)
        //add space
        var space = document.createElement("br");
        inside2.appendChild(space);
        inside2.appendChild(space);
        inside2.appendChild(form);
        inside2.appendChild(form2);
        if(obj.is_admin == "yes"){
            inside2.appendChild(form3);
            inside2.appendChild(space);
        }
        inside.appendChild(inside2);
        body.appendChild(inside);
        
        
    }
}) 

//document.getElementById("body").createElement('footer').appendChild()