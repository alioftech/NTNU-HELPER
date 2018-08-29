with Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;
use  Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;

procedure exercise7 is

    Count_Failed    : exception;    -- Exception to be raised when counting fails
    Gen             : Generator;    -- Random number generator

    protected type Transaction_Manager (N : Positive) is
        entry Finished;
        function Commit return Boolean;
        procedure Signal_Abort;
    private
        Finished_Gate_Open  : Boolean := False;
        Aborted             : Boolean := False;
        Should_Commit       : Boolean := True;
    end Transaction_Manager;
    protected body Transaction_Manager is
        entry Finished when Finished_Gate_Open or Finished'Count = N is
        begin
            ------------------------------------------
            -- PART 3: Complete the exit protocol here
            Should_Commit := not Aborted;
            Finished_Gate_Open := Finished'Count /= 0;
            if not Finished_Gate_Open then
                Aborted := False;
            end if;
            ------------------------------------------
        end Finished;

        procedure Signal_Abort is
        begin
            Aborted := True;
        end Signal_Abort;

        function Commit return Boolean is
        begin
            return Should_Commit;
        end Commit;
        
    end Transaction_Manager;



    
    function Unreliable_Slow_Add (x : Integer) return Integer is
    Error_Rate : Constant := 0.15;  -- (between 0 and 1)
    begin
        -------------------------------------------
        -- PART 1: Create the transaction work here
        if Random(Gen) > Error_Rate then
            delay Duration(3.5 + Random(Gen));
            return x + 10;
        else
            delay Duration(0.5 * Random(Gen));
            raise Count_Failed;
        end if;
        -------------------------------------------
    end Unreliable_Slow_Add;




    task type Transaction_Worker (Initial : Integer; Manager : access Transaction_Manager);
    task body Transaction_Worker is
        Num         : Integer   := Initial;
        Prev        : Integer   := Num;
        Round_Num   : Integer   := 0;
    begin
        Put_Line ("Worker" & Integer'Image(Initial) & " started");

        loop
            Put_Line ("Worker" & Integer'Image(Initial) & " started round" & Integer'Image(Round_Num));
            Round_Num := Round_Num + 1;

            -------------------------------------------
            -- PART 2: Do the transaction work here
            begin
                Num := Unreliable_Slow_Add(Prev);
            exception
                when Count_Failed =>
                    Manager.Signal_Abort;
            end;
            Manager.Finished;
            -------------------------------------------

            if Manager.Commit = True then
                Put_Line ("  Worker" & Integer'Image(Initial) & " comitting" & Integer'Image(Num));
            else
                Put_Line ("  Worker" & Integer'Image(Initial) &
                             " reverting from" & Integer'Image(Num) &
                             " to" & Integer'Image(Prev));
                -------------------------------------------
                -- PART 2: Roll back to previous value here
                Num := Prev;
                -------------------------------------------
            end if;

            Prev := Num;
            delay 0.5;

        end loop;
    end Transaction_Worker;

    Manager : aliased Transaction_Manager (3);

    Worker_1 : Transaction_Worker (0, Manager'Access);
    Worker_2 : Transaction_Worker (1, Manager'Access);
    Worker_3 : Transaction_Worker (2, Manager'Access);

begin
    Reset(Gen); -- Seed the random number generator
end exercise7;
