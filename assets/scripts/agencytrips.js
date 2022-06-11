fetch("/jsonAllTrips/")
.then(response => response.json())
.then(data => {
    var body2 = document.getElementById("body");
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        if(obj.agencyName == document.getElementById("theUsername").innerHTML){
        //add the text
        const paragraph = document.createElement("h2");
        paragraph.setAttribute("id","agency-paragraph"+i);
        paragraph.innerHTML="Agentia de turism "+obj.agencyName;
        //add the text
        const paragraph2 = document.createElement("p");
        paragraph2.setAttribute("id","title-paragraph"+i);

        paragraph2.innerHTML=obj.title;
        //add the text
        const paragraph3 = document.createElement("p");
        paragraph3.setAttribute("id","desc-paragraph"+i);
    
        paragraph3.innerHTML=obj.description;
        //add the text
        const paragraph7 = document.createElement("p");
        paragraph7.setAttribute("id","location-paragraph"+i);
   
        paragraph7.innerHTML="Orasul: " + obj.city + " Tara: " + obj.country;
        //add the text
        const paragraph6 = document.createElement("p");
        paragraph6.setAttribute("id","hotel-paragraph"+i);
        
        paragraph6.innerHTML="Hotelul " + obj.hotel +" cu "+obj.stars+" stele";
        //add the text
        const paragraph4 = document.createElement("p");
        paragraph4.setAttribute("id","price-paragraph"+i);
        
        paragraph4.innerHTML=obj.price+" euro";
        //add the text
        const paragraph5 = document.createElement("p");
        paragraph5.setAttribute("id","details-paragraph"+i);
  
        paragraph5.innerHTML="De la data de "+ obj.date+" timp de "+obj.days+" zile";
        var divGalerie = document.createElement("div");
        divGalerie.setAttribute("id","grid-galerie");
        //add the photo
        var img1 = document.createElement("img");
        img1.setAttribute("id","img1"+i);   
        img1.src = "../"+obj.img1     

        img1.setAttribute("class","card-img-top")    
        img1.setAttribute("width","100%")   
        img1.setAttribute("height","100%")   
      //  child.appendChild(img1)
        //add the photo
        var img2 = document.createElement("img");
        img2.setAttribute("id","img2"+i);        
        img2.src = "../"+obj.img2
        img2.setAttribute("class","card-img-top")    
        img2.setAttribute("width","100%")   
        img2.setAttribute("height","100%")   
        //child.appendChild(img2)
        //add the photo
        var img3 = document.createElement("img");
        img3.setAttribute("id","img3"+i);        
        img3.src = "../"+obj.img3
        img3.setAttribute("class","card-img-top")    
        img3.setAttribute("width","100%")   
        img3.setAttribute("height","100%")   
        //child.appendChild(img3)
        divGalerie.appendChild(img1);
        divGalerie.appendChild(img2);
        divGalerie.appendChild(img3);
  
        //add space
        var space = document.createElement("br");
        var separator = document.createElement("hr");
        


        //add a button
        if(obj.same == "yes"){
            var form1 = document.createElement("form");
            form1.setAttribute("id","deleteTrip");
            form1.method="GET";
            form1.action="/deleteTrip/"+obj.id;
            
            var button = document.createElement("button");
            button.innerHTML = "Sterge excursia";
            button.setAttribute("id","delte-button");
            button.setAttribute("class","toRegister-btn")
            button.type="submit";
            form1.appendChild(button); 
            

            var form2 = document.createElement("form");
            form2.setAttribute("id","updateTrip");
            form2.method="GET";
            form2.action="/updateTripPage/"+obj.id;
            
            var button2 = document.createElement("button");
            button2.innerHTML = "Editeaza excursia";
            button2.setAttribute("id","update-button");
            button2.setAttribute("class","toRegister-btn")
            button2.type="submit";
            form2.appendChild(button2); 
            
            }

            var form3 = document.createElement("form");
            form3.setAttribute("id","getWeather");
            form3.method="GET";
            form3.action="/weather/"+obj.city;  
            
            var button3 = document.createElement("button");
            button3.innerHTML = "Vremea in "+ obj.city;
            button3.type="submit";
            button3.setAttribute("id","weather");
            button3.setAttribute("class","toRegister-btn")
            form3.appendChild(button3); 

            var body = document.createElement("div");
            body.setAttribute("class","formular-login2");
            body.appendChild(paragraph);
            body.appendChild(paragraph2);
            body.appendChild(paragraph3);
            body.appendChild(paragraph7);
            body.appendChild(paragraph6);
            body.appendChild(paragraph4);
            body.appendChild(paragraph5);
            body.appendChild(divGalerie);
            body.appendChild(space);
            if(obj.same == "yes"){
            
            body.appendChild(form1);
            body.append(space);
            body.appendChild(form2);
            body.append(space);
            }
            
            body.appendChild(form3);
            body2.appendChild(body);
        }
        
       

    }
}) 