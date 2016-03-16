$(document).ready(function(){
	// login to server
    $("#login").click(function(){
    	$("#join").prop('disabled', false);
        $.get("/login?rcpt=" + from.value, function(data, status){
        	if(data == "user reg"){
        		alert("name have been used");
           	}else{
           		$("#rooms").prop('disabled', false);
  				$("#login").prop('disabled', true);
  				$("#from").prop('disabled', true);
  				$("#logout").prop('disabled', false);
				var res = data.split(",");
           		for (i = 0; i < res.length; i++) { 
    				if(res[i] != ""){
    					$('#rooms').append('<option value="5">'+ res[i] + '</option>');
    				}
				}
				$("#rooms").val($("#rooms option:first").val());
				longpolluser("/joinroom?from=" + from.value, recv);
           	}
    	});
    });
    
    // join room and send notification
    $("#join").click(function(){
    	
    	$.get("/notify?from=" + from.value + "&room=" + $("#rooms option:selected").text(), function(data, status){
        	
        	if(data == "success"){
        		
           	}else{
           		
           	}
    	});
    	
    });
    
    
    //leave room and send notification
    $("#leave").click(function(){
    	
    	$.get("/leaveroom?from=" + from.value + "&room=" + $("#rooms option:selected").text(), function(data, status){
        	
        	if(data == "success"){
        		
           	}else{
           		
           	}
    	});
    	$("#users").empty();
        $("#users").prop('disabled', true);
        $("#rooms").prop('disabled', false);
        $("#leave").prop('disabled', true);
        $("#join").prop('disabled', false);
    });
    // logout from server
    $("#logout").click(function(){
    	$( "#leave" ).trigger( "click" );
    	$.get("/logout?from=" + from.value, function(data, status){
        	
        	if(data == "success"){
        		$("#login").prop('disabled', false);
        		$("#from").prop('disabled', false);
        		$("#rooms").empty();
        		$("#rooms").prop('disabled', true);
        		$("#logout").prop('disabled', true);
        		
           	}else{
           		alert("No this user");
           	}
    	});
    	
    });

    
    // select one user to send information
    $("#users")
  	.change(function() {
    var str = "";
    $( "select option:selected" ).each(function() {
      	str = $(this).text();
    });	
    	$("#rcpt").val(str);
  	})
  	.trigger( "change" );
  	
  	$("#rcpt").prop('disabled', true);
  	$("#rooms").prop('disabled', true);
  	$("#users").prop('disabled', true);
  	$("#join").prop('disabled', true);
  	$("#leave").prop('disabled', true);
  	$("#logout").prop('disabled', true);
});
//send information
function send() {

	var from = document.getElementById("from");
	var rcpt = document.getElementById("rcpt");
	var box = document.getElementById("box");
	var input = document.getElementById("input");

	var req = new XMLHttpRequest (); 
	req.open ('POST', "/notify?from="+from.value+"&rcpt=" + rcpt.value + "&room="+$( "#rooms option:selected" ).text(), true); 

	req.onreadystatechange = function (aEvt) {
		if (req.readyState == 4) { 
			if (req.status == 200) {
				//alert("send");
			} else {
				//alert ("failed to send!");
			}
		}
	};
	req.send(input.value);
	box.value += "\nme: " + input.value;
	input.value = "";
	
}
// poll chat information or user list to client
function longpolluser(url, callback) {
	$.get(url, function(data, status){
        	
        	if(data == "No this room"){
        		alert("No this room");
        		return;
           	}else if(data == "leave"){
           		return;
           	}
           	else{
           		
           		$("#users").empty();
				var res = data.split(",");
				$('#users').append('<option value="5">All</option>');
           		for (i = 0; i < res.length-2; i++) { 
    				if(res[i] != ""){
    					$('#users').append('<option value="5">'+ res[i] + '</option>');
    				}
				}
				$("#users").val($("#rooms option:first").val());
				$("#users").prop('disabled', false);
        		$("#rooms").prop('disabled', true);
        		$("#leave").prop('disabled', false);
        		$("#join").prop('disabled', true);
        		if($("#from").val() != res[res.length-2]){
        			
        			callback(res[res.length-2]+" "+res[res.length-1]);
        		}
        		if(res[res.length-1] == "leave" && res[res.length-2] == $("#from").val()){
        			return;
        		}
        		longpolluser(url, callback);
           	}
     });
	req.send(null);
}
//show information
function recv(msg) {
	var box = document.getElementById("box");
	box.value += "\n" + msg;
}