function viewGame(){
        $("div.game-page ul.ui-listview").html("");
        $.mobile.changePage("#pageGame");
}

function viewResult(gameId){
    $.mobile.changePage("#pageGameResult?id="+gameId);
}

function loadGameResult(gameId){
      console.log("init game result:" + gameId);
      var $table = $('#game_result_table');

      $table.bootstrapTable('destroy');
      CG.PostController.get('/v1/games/'+gameId, function(err, data){
          var name = data.GameTemplate.Name;
          var var1Type = data.GameTemplate.Var1Type;
          var var1Selects = [];
          var type = data.GameTemplate.Type;
          var subtype = data.GameTemplate.Subtype;

          if (var1Type == 3){
              var s = data.GameTemplate.Var1Select.split(",");
              for (var i = 0; i < s.length; i++){
                  var1Selects[s[i].split(":")[0]] = s[i].split(":")[1];
              }
          }
          if (data.GameTemplate.Subname){
             name += " - "+data.GameTemplate.Subname;
          }

          if (data.ShowResult == 1){
              $("#result_toolbar").show();
              $('#game_result_table').show();
          }
          else{
              $("#result_toolbar").hide();
              $('#game_result_table').hide();
          }

          $("#result_game_name").text(filterXSS(name+" 游戏结果"));

          CG.PostController.get('/v1/games/'+gameId+'/stat', function(err, data){
            if (type == 2 && subtype == 3){
                for (var i = 0; i < data.length; i++){
                    data[i].Var1 = data[i].Var1 +" - "+ data[i].Var2;
                }
            }
            else if (var1Type == 3){
                for (var i = 0; i < data.length; i++){
                    data[i].Var1 = var1Selects[data[i].Var1];
                }
            }
            console.log(data);
            $table.bootstrapTable({data: data, 
              columns: [{
                field: 'Var1',
                title: '数字'
              },
              {
                field: 'Count',
                title: '人数'
              }]
            });
          });
          

          var $win = $('#game_win_table');
          $win.bootstrapTable('destroy');
          CG.PostController.get('/v1/games/'+gameId+'/win', function(err, data){
              if (data){
                if (var1Type == 3){
                    for (var i = 0; i < data.length; i++){
                        data[i].Var1 = var1Selects[data[i].Var1];
                    }
                }
                for (var i = 0; i < data.length; i++){
                    var name = "";
                    if (data[i].User.Fullname){
                        name += data[i].User.Fullname;
                    }
                    if (data[i].IsWin == 2){
                        name += "(获得书)"
                    }
                    else{
                        name += "(获胜)"
                    }
                    data[i].Name = name;
                    data[i].StudentNo = data[i].User.StudentNo;
                }
                $win.bootstrapTable({data: data});
             }
             else{
                $win.bootstrapTable({data: []});
             }
          });
          
      });

      
}

function refreshGame(){
    loadGameResult($('#result_game_id').val());
}

function activateGame(){
    //{"templateId":1, "reward":5, "gameTime":120, "playerNum":10,  "showResult":1}
    var templateId = $('#game_template_id').val();
    if (templateId == 0){
        templateId = $('#game_subtype').val();
    }

    var postId = window.sessionStorage.getItem("postId");
    var data = {"templateId":templateId, "reward":$('#game_reward').val(), "gameTime":$('#game_time').val(), "playerNum":$('#game_playerNum').val(), "showResult":$('#game_showResult').val()};
    CG.PostController.post('/v1/posts/' + postId + '/game', data, function(err, data){
        if (!err){
            viewGame();
        }
    });
}

function delayHideError(){
   setTimeout(function(){
      $("#game_error_div").hide();
   }, 5000);
}

function submitPlayGame(){
   var gameId = $("#player_game_id").val();

   if (parseInt($('#player_game_var1').val()) < 0 || parseInt($('#player_game_var1').val()) > 100){
      $("#game_error_msg").text("请输入0-100整数");
      $("#game_error_div").show();
      delayHideError();
      return;
   }

   if (parseInt($('#player_game_var2').val()) < 0 || parseInt($('#player_game_var2').val()) > 100){
      $("#game_error_msg").text("请输入0-100整数");
      $("#game_error_div").show();
      delayHideError();
      return;
   }

   var data = {"var1":$('#player_game_var1').val(), "var2":$('#player_game_var2').val()};
   if ($("#player_game_var3_field").attr("style") === "display: block;"){
      data.var1 = $('#player_game_var3_select').val();
   } 

   CG.PostController.post('/v1/games/' + gameId , data, function(err, data){
      if (!err){
          viewResult(gameId);
      }
   });
}

function newGameDetail(id){   
    CG.PostController.get('/v1/games/template/'+id, function(err, data){
        var size = data.length;
        if (size > 0){
            var game1 = data[0];

            $('#game_name').text(filterXSS(game1.Name));
            $('#game_type').text(filterXSS(""+game1.Type));
            $('#game_reward').val("");
            $('#game_playerNum').val("");

            if (game1.Type == 2  || game1.Type == 4 || game1.Type == 5 || game1.Type == 6){
                $('#game_subtype_field').show();
                $('#game_template_id').val(0);
            }
            else{
                $('#game_subtype_field').hide();
                $('#game_template_id').val(id);
            }
            
            if (game1.Type == 1){
                $('#game_playerNum_field').show();
            }
            else{
                $('#game_playerNum_field').hide();
            }

            if (game1.Type == 1 || game1.Type == 3 || game1.Type == 4 || game1.Type == 5 || game1.Type == 6){
                $('#game_reward_field').show();
            }
            else{
                $('#game_reward_field').hide();
            }
        }

        if (size > 1){
            $('#game_subtype').html("");
            $("#game_subtype-button span").html("<span>&nbsp;"+filterXSS(data[0].Subname)+"</span>");
            for (var i = 0; i < size; i++){
                $('#game_subtype').append('<option value="'+ data[i].Id+'" >'+filterXSS(data[i].Subname)+'</option>');
            }
             
        }
    });
}

function countdown(time){
    if (time <= 0){
        $("#player_game_time").text("0秒");
        return;
    }

    var t = setInterval(function(){
      console.log("ontime"); 
      time--;
      $("#player_game_time").text(time + "秒");
      if (time == 0){
        clearInterval(t);
      }
    }, 1000);
}

function viewActivateGame(id, status) {
  CG.PostController.get('/v1/games/'+id, function(err, data){
      if (!err){
          var name = data.name;
          if (data.subname){
             name += " - "+data.subname;
          }
          $("#player_game_var1").val("");
          $("#player_game_var2").val("");
          // if teacher or game is over redirect to result
          if (checkCurrentUserId(data.userId) || data.status != 1){
             viewResult(id);
          }
          else{
             // title
             $("#player_game_name").text(filterXSS(name));
             $("#player_game_id").val(id);
             $("#player_game_time").text(filterXSS(data.restTime + "秒"));
             countdown(data.restTime);
             
             
             // rule
             $("#player_game_rule").text(filterXSS(data.ruleLabel));
             $("#player_game_var1help").text(filterXSS(data.var1Help));
             $("#player_game_var1").attr("placeholder", data.var1Help);
             $("#player_game_var2help").text(filterXSS(data.var2Help));
             $("#player_game_var2").attr("placeholder", data.var2Help);
             if(data.var1Label){
                $("#player_game_var1label").text(filterXSS(data.var1Label));
             }
             else{
                $("#player_game_var1label").text("");
             }
             if(data.var2Label){
                $("#player_game_var2label").text(filterXSS(data.var2Label));
             }
             else{
                $("#player_game_var2label").text("");
             }

             $("#player_game_var1_field").show();
             $("#player_game_var3_field").hide();

             if ((data.type == 2 && data.subtype == 3) ){
                $("#player_game_var2_field").show();
                $("#player_game_var1help_field").hide();
                $("#player_game_var2help_field").hide();
             }
             else if (data.type == 3){
                $("#player_game_var1help_field").show();
                $("#player_game_var2help_field").hide();
                $("#player_game_var2_field").hide();
             }
             else if ((data.type == 4 && data.subtype == 3) || (data.type == 6 && data.subtype == 2) || (data.type == 5 && data.subtype == 2) || (data.type == 5 && data.subtype == 3)){
                $("#player_game_var1help_field").hide();
                $("#player_game_var2help_field").hide();
                $("#player_game_var1_field").hide();
                $("#player_game_var2_field").hide();
                $("#player_game_var3_field").show();

                $('#player_game_var3_select').html("");
                var addrs = data.var1Select.split(",");
                $("#player_game_var3_select span").html("<span>&nbsp;"+filterXSS(addrs[0].split(":")[1])+"</span>");
                $('#player_game_var3_select').val("0");
                for (var i = 0; i < addrs.length; i++){
                  $('#player_game_var3_select').append('<option value="'+ addrs[i].split(":")[0]+'" >'+filterXSS(addrs[i].split(":")[1])+'</option>');
                }
             }
             else{
                $("#player_game_var1help_field").hide();
                $("#player_game_var2_field").hide();
                $("#player_game_var2help_field").hide();
             }
             $.mobile.changePage("#pagePlayer?id="+id);
          }
      } 
  });
}

// pageGame
(function pullPagePullImplementation($) {
  "use strict";
  

   function getGameTemplateData(skip, callback){
        CG.PostController.get('/v1/games/template', function(err, data){
            callback(err, data);
        });
    }
    
    function initGameTemplate(){
        var listSelector = "div.newgamelist-page ul.ui-listview";
        if ($(listSelector) && $(listSelector).html() && $(listSelector).html().length > 0){
            ;
        }
        else {
            getGameTemplateData(0, function(err, content){
                if (!err){
                    var i, newContent = "";
                    content.forEach(function(e){
                        var html = '<li data-icon="false"><a href="#pageNewGame?id='+ e.Id +'">';
                        html += '<h2>' + filterXSS(e.Name) + '</h2>';
                        html += '</a></li>'
                        newContent = html + newContent;
                }); 

                $(listSelector).prepend(newContent).listview("refresh");  // Prepend new content and refresh listview                      
            }});            
        }
    }

    function initGameData() {
      CG.GamePullDownRefresh.initData();
    }

    function initGameResult(gameId){
      $('#result_game_id').val(gameId);
      loadGameResult(gameId);
    }
    
  // This is the callback that is called when the user has completed the pull-down gesture.
  // Your code should initiate retrieving data from a server, local database, etc.
  // Typically, you will call some function like jQuery.ajax() that will make a callback
  // once data has been retrieved.
  //
  // For demo, we just use timeout to simulate the time required to complete the operation.
  function onPullDown (event, data) {
      CG.GamePullDownRefresh.onPullDown(event, data);
  }    

  // Called when the user completes the pull-up gesture.
  function onPullUp (event, data) { 
      CG.GamePullDownRefresh.onPullUp(event, data);
  }    
  
  // Set-up jQuery event callbacks
  $(document).delegate("div.game-page", "pageinit", 
    function bindPullPagePullCallbacks(event) {
      $(".iscroll-wrapper", this).bind( {
      iscroll_onpulldown : onPullDown,
      iscroll_onpullup   : onPullUp
      } );
    } );  

  $(document).on("pagebeforechange", function(e, f){
    if (typeof f.toPage !== "string"){
        return;
    }

    //var paramUrl = $.mobile.path.parseUrl(f.toPage);
    //console.log(paramUrl);

    var hashs = $.mobile.path.parseUrl(f.absUrl).hash.split("?");
    var hash = hashs[0];
    var search = hashs[1];

    if (hash === "#pageGame"){
        initGameData();
    }
    else if (hash === "#pageNewGame"){
      if (search){
        var id = search.split("=")[1];
        newGameDetail(id);
      }
    }
    else if (hash === "#pageNewGamelist"){
       initGameTemplate();
    }
    else if (hash === "#pageViewGame"){
      if (search){
        var id = search.split("=")[1];
        viewGame(id);
      }
    }
    else if (hash === "#pageGameResult"){
      if (search){
        var id = search.split("=")[1];
        initGameResult(id);
      }
    }

  });

}(jQuery));

